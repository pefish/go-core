package api_strategy

import (
	"github.com/golang/mock/gomock"
	mock_api_session "github.com/pefish/go-core/mock/mock-api-session"
	go_error "github.com/pefish/go-error"
	go_jwt "github.com/pefish/go-jwt"
	"github.com/pefish/go-test-assert"
	"testing"
	"time"
)

func TestJwtAuthStrategyClass_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	apiSessionInstance := mock_api_session.NewMockIApiSession(ctrl)
	apiSessionInstance.EXPECT().Header("jwt").Return("gsfdg")
	var userIdResult uint64
	apiSessionInstance.EXPECT().SetJwtBody(gomock.Any()).AnyTimes()
	apiSessionInstance.EXPECT().SetUserId(gomock.Any()).DoAndReturn(func(userId uint64) {
		userIdResult = userId
	}).AnyTimes()
	apiSessionInstance.EXPECT().SetJwtHeaderName(gomock.Any()).AnyTimes()
	JwtAuthApiStrategy.SetHeaderName("jwt")
	err := JwtAuthApiStrategy.Execute(apiSessionInstance, nil)
	test.Equal(t, JwtAuthApiStrategy.GetErrorCode(), err.ErrorCode)
	test.Equal(t, "Unauthorized", err.ErrorMessage)


	pkey, pubkey, err1 := go_jwt.GeneRsaKeyPair()
	test.Equal(t, nil, err1)
	jwt, err2 := go_jwt.Jwt.GetJwt(pkey, 60 * time.Second, map[string]interface{}{
		"user_id": 6356,
	})
	test.Equal(t, nil, err2)
	JwtAuthApiStrategy.SetPubKey(pubkey)
	apiSessionInstance.EXPECT().Header("jwt").Return(jwt).AnyTimes()
	apiSessionInstance.EXPECT().Data(gomock.Any()).AnyTimes()
	apiSessionInstance.EXPECT().SetData(gomock.Any(), gomock.Any()).AnyTimes()
	err = JwtAuthApiStrategy.Execute(apiSessionInstance, nil)
	test.Equal(t, (*go_error.ErrorInfo)(nil), err)
	test.Equal(t, uint64(6356), userIdResult)
}