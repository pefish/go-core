package global_api_strategy

import (
	"context"
	"contrib.go.opencensus.io/exporter/stackdriver"
	go_application "github.com/pefish/go-application"
	"github.com/pefish/go-core/api-session"
	"github.com/pefish/go-core/driver/logger"
	go_error "github.com/pefish/go-error"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type OpenCensusClass struct {
	errorCode uint64
}

var OpenCensusStrategy = OpenCensusClass{}

func (this *OpenCensusClass) GetName() string {
	return `OpenCensus`
}

func (this *OpenCensusClass) GetDescription() string {
	return `OpenCensus`
}

func (this *OpenCensusClass) SetErrorCode(code uint64) {
	this.errorCode = code
}

func (this *OpenCensusClass) GetErrorCode() uint64 {
	if this.errorCode == 0 {
		return go_error.INTERNAL_ERROR_CODE
	}
	return this.errorCode
}

type OpenCensusStrategyParam struct {
	StackDriverOption *stackdriver.Options
	EnableTrace       bool
	EnableStats       bool
}

func (this *OpenCensusClass) Init(param interface{}) {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s Init`, this.GetName())
	defer logger.LoggerDriver.Logger.DebugF(`api-strategy %s Init defer`, this.GetName())
	if param == nil {
		go_error.Throw(`OpenCensusStrategyParam must be set`, this.GetErrorCode())
	}
	go func() {
		option := stackdriver.Options{
			ReportingInterval: 60 * time.Second,
		}
		newParam := param.(OpenCensusStrategyParam)
		if newParam.StackDriverOption != nil {
			option.ProjectID = newParam.StackDriverOption.ProjectID
		}
		exporter, err := stackdriver.NewExporter(option)
		if err != nil {
			panic(err)
		}
		defer exporter.Flush()
		if newParam.EnableStats {
			err = exporter.StartMetricsExporter()
			if err != nil {
				panic(err)
			}
			defer exporter.StopMetricsExporter()
		}
		trace.RegisterExporter(exporter)
		defer trace.UnregisterExporter(exporter)
		if go_application.Application.Env == `local` { // 本地调试才打开
			trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()}) // 每个请求一个trace，生产环境不要使用
		}
		select {
		case <- go_application.Application.OnFinished():
			break
		}
	}()
}

func (this *OpenCensusClass) Execute(out *api_session.ApiSessionClass, param interface{}) {
	logger.LoggerDriver.Logger.DebugF(`api-strategy %s trigger`, this.GetName())
	defer func() {
		if err := recover(); err != nil {
			logger.LoggerDriver.Logger.Error(err)
		}
	}()
	newParam := param.(OpenCensusStrategyParam)
	w, r := out.ResponseWriter, out.Request
	if newParam.EnableTrace {
		r1, traceEnd := startTrace(w, r)
		out.AddDefer(func() {
			traceEnd()
		})
		r = r1
	}
	if newParam.EnableStats {
		var tags addedTags
		_, statsEnd := startStats(w, r)
		out.AddDefer(func() {
			statsEnd(&tags)
		})
	}
}

// -----------------

var defaultFormat propagation.HTTPFormat = &b3.HTTPFormat{}

func startTrace(w http.ResponseWriter, r *http.Request) (*http.Request, func()) {
	name := r.URL.Path
	ctx := r.Context()

	startOpts := trace.StartOptions{}
	var span *trace.Span
	sc, ok := defaultFormat.SpanContextFromRequest(r)
	if ok {
		ctx, span = trace.StartSpanWithRemoteParent(ctx, name, sc,
			trace.WithSampler(startOpts.Sampler),
			trace.WithSpanKind(trace.SpanKindServer))
	} else {
		ctx, span = trace.StartSpan(ctx, name,
			trace.WithSampler(startOpts.Sampler),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		if ok {
			span.AddLink(trace.Link{
				TraceID:    sc.TraceID,
				SpanID:     sc.SpanID,
				Type:       trace.LinkTypeParent,
				Attributes: nil,
			})
		}
	}
	span.AddAttributes(requestAttrs(r)...)
	return r.WithContext(ctx), span.End
}

func requestAttrs(r *http.Request) []trace.Attribute {
	userAgent := r.UserAgent()

	attrs := make([]trace.Attribute, 0, 5)
	attrs = append(attrs,
		trace.StringAttribute(ochttp.PathAttribute, r.URL.Path),
		trace.StringAttribute(ochttp.URLAttribute, r.URL.String()),
		trace.StringAttribute(ochttp.HostAttribute, r.Host),
		trace.StringAttribute(ochttp.MethodAttribute, r.Method),
	)

	if userAgent != "" {
		attrs = append(attrs, trace.StringAttribute(ochttp.UserAgentAttribute, userAgent))
	}

	return attrs
}

func startStats(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, func(tags *addedTags)) {
	ctx, _ := tag.New(r.Context(),
		tag.Upsert(ochttp.Host, r.Host),
		tag.Upsert(ochttp.Path, r.URL.Path),
		tag.Upsert(ochttp.Method, r.Method))
	track := &trackingResponseWriter{
		start:  time.Now(),
		ctx:    ctx,
		writer: w,
	}
	if r.Body == nil {
		// TODO: Handle cases where ContentLength is not set.
		track.reqSize = -1
	} else if r.ContentLength > 0 {
		track.reqSize = r.ContentLength
	}
	stats.Record(ctx, ochttp.ServerRequestCount.M(1))
	return track.wrappedResponseWriter(), track.end
}

type addedTags struct {
	t []tag.Mutator
}

type trackingResponseWriter struct {
	ctx        context.Context
	reqSize    int64
	respSize   int64
	start      time.Time
	statusCode int
	statusLine string
	endOnce    sync.Once
	writer     http.ResponseWriter
}

func (t *trackingResponseWriter) end(tags *addedTags) {
	t.endOnce.Do(func() {
		if t.statusCode == 0 {
			t.statusCode = 200
		}

		span := trace.FromContext(t.ctx)
		span.SetStatus(ochttp.TraceStatus(t.statusCode, t.statusLine))
		span.AddAttributes(trace.Int64Attribute(ochttp.StatusCodeAttribute, int64(t.statusCode)))

		m := []stats.Measurement{
			ochttp.ServerLatency.M(float64(time.Since(t.start)) / float64(time.Millisecond)),
			ochttp.ServerResponseBytes.M(t.respSize),
		}
		if t.reqSize >= 0 {
			m = append(m, ochttp.ServerRequestBytes.M(t.reqSize))
		}
		allTags := make([]tag.Mutator, len(tags.t)+1)
		allTags[0] = tag.Upsert(ochttp.StatusCode, strconv.Itoa(t.statusCode))
		copy(allTags[1:], tags.t)
		stats.RecordWithTags(t.ctx, allTags, m...)
	})
}

func (t *trackingResponseWriter) Header() http.Header {
	return t.writer.Header()
}

func (t *trackingResponseWriter) Write(data []byte) (int, error) {
	n, err := t.writer.Write(data)
	t.respSize += int64(n)
	return n, err
}

func (t *trackingResponseWriter) WriteHeader(statusCode int) {
	t.writer.WriteHeader(statusCode)
	t.statusCode = statusCode
	t.statusLine = http.StatusText(t.statusCode)
}

func (t *trackingResponseWriter) wrappedResponseWriter() http.ResponseWriter {
	var (
		hj, i0 = t.writer.(http.Hijacker)
		cn, i1 = t.writer.(http.CloseNotifier)
		pu, i2 = t.writer.(http.Pusher)
		fl, i3 = t.writer.(http.Flusher)
		rf, i4 = t.writer.(io.ReaderFrom)
	)

	switch {
	case !i0 && !i1 && !i2 && !i3 && !i4:
		return struct {
			http.ResponseWriter
		}{t}
	case !i0 && !i1 && !i2 && !i3 && i4:
		return struct {
			http.ResponseWriter
			io.ReaderFrom
		}{t, rf}
	case !i0 && !i1 && !i2 && i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Flusher
		}{t, fl}
	case !i0 && !i1 && !i2 && i3 && i4:
		return struct {
			http.ResponseWriter
			http.Flusher
			io.ReaderFrom
		}{t, fl, rf}
	case !i0 && !i1 && i2 && !i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Pusher
		}{t, pu}
	case !i0 && !i1 && i2 && !i3 && i4:
		return struct {
			http.ResponseWriter
			http.Pusher
			io.ReaderFrom
		}{t, pu, rf}
	case !i0 && !i1 && i2 && i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Pusher
			http.Flusher
		}{t, pu, fl}
	case !i0 && !i1 && i2 && i3 && i4:
		return struct {
			http.ResponseWriter
			http.Pusher
			http.Flusher
			io.ReaderFrom
		}{t, pu, fl, rf}
	case !i0 && i1 && !i2 && !i3 && !i4:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
		}{t, cn}
	case !i0 && i1 && !i2 && !i3 && i4:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			io.ReaderFrom
		}{t, cn, rf}
	case !i0 && i1 && !i2 && i3 && !i4:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Flusher
		}{t, cn, fl}
	case !i0 && i1 && !i2 && i3 && i4:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Flusher
			io.ReaderFrom
		}{t, cn, fl, rf}
	case !i0 && i1 && i2 && !i3 && !i4:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Pusher
		}{t, cn, pu}
	case !i0 && i1 && i2 && !i3 && i4:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Pusher
			io.ReaderFrom
		}{t, cn, pu, rf}
	case !i0 && i1 && i2 && i3 && !i4:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Pusher
			http.Flusher
		}{t, cn, pu, fl}
	case !i0 && i1 && i2 && i3 && i4:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Pusher
			http.Flusher
			io.ReaderFrom
		}{t, cn, pu, fl, rf}
	case i0 && !i1 && !i2 && !i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
		}{t, hj}
	case i0 && !i1 && !i2 && !i3 && i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			io.ReaderFrom
		}{t, hj, rf}
	case i0 && !i1 && !i2 && i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.Flusher
		}{t, hj, fl}
	case i0 && !i1 && !i2 && i3 && i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.Flusher
			io.ReaderFrom
		}{t, hj, fl, rf}
	case i0 && !i1 && i2 && !i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.Pusher
		}{t, hj, pu}
	case i0 && !i1 && i2 && !i3 && i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.Pusher
			io.ReaderFrom
		}{t, hj, pu, rf}
	case i0 && !i1 && i2 && i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.Pusher
			http.Flusher
		}{t, hj, pu, fl}
	case i0 && !i1 && i2 && i3 && i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.Pusher
			http.Flusher
			io.ReaderFrom
		}{t, hj, pu, fl, rf}
	case i0 && i1 && !i2 && !i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.CloseNotifier
		}{t, hj, cn}
	case i0 && i1 && !i2 && !i3 && i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.CloseNotifier
			io.ReaderFrom
		}{t, hj, cn, rf}
	case i0 && i1 && !i2 && i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.CloseNotifier
			http.Flusher
		}{t, hj, cn, fl}
	case i0 && i1 && !i2 && i3 && i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.CloseNotifier
			http.Flusher
			io.ReaderFrom
		}{t, hj, cn, fl, rf}
	case i0 && i1 && i2 && !i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.CloseNotifier
			http.Pusher
		}{t, hj, cn, pu}
	case i0 && i1 && i2 && !i3 && i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.CloseNotifier
			http.Pusher
			io.ReaderFrom
		}{t, hj, cn, pu, rf}
	case i0 && i1 && i2 && i3 && !i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.CloseNotifier
			http.Pusher
			http.Flusher
		}{t, hj, cn, pu, fl}
	case i0 && i1 && i2 && i3 && i4:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.CloseNotifier
			http.Pusher
			http.Flusher
			io.ReaderFrom
		}{t, hj, cn, pu, fl, rf}
	default:
		return struct {
			http.ResponseWriter
		}{t}
	}
}
