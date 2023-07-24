// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package server

import (
	pb "cloudsave-validator-grpc-plugin-server-go/pkg/pb"
	"context"
	"encoding/json"
	"strings"
	"time"
)

type CloudsaveValidatorServer struct {
	pb.UnimplementedCloudsaveValidatorServiceServer
}

func (s *CloudsaveValidatorServer) BeforeWriteGameRecord(ctx context.Context, request *pb.GameRecord) (*pb.GameRecordValidationResult, error) {
	if strings.HasSuffix(request.GetKey(), "map") {
		var r CustomGameRecord
		err := json.Unmarshal(request.GetPayload(), &r)
		if err != nil {
			return nil, err
		}
		if err = r.Validate(); err != nil {
			errorDetail := &pb.Error{
				ErrorCode:    1,
				ErrorMessage: err.Error(),
			}

			return &pb.GameRecordValidationResult{
				IsSuccess: false,
				Key:       request.Key,
				Error:     errorDetail,
			}, nil
		}
	}

	return &pb.GameRecordValidationResult{IsSuccess: true, Key: request.Key}, nil
}

func (s *CloudsaveValidatorServer) AfterReadGameRecord(ctx context.Context, gameRecord *pb.GameRecord) (*pb.GameRecordValidationResult, error) {
	if strings.HasSuffix(gameRecord.Key, "daily_msg") {
		var r DailyMessage
		err := json.Unmarshal(gameRecord.GetPayload(), &r)
		if err != nil {
			return nil, err
		}
		if time.Now().Before(r.AvailableOn) {
			return &pb.GameRecordValidationResult{
				IsSuccess: false,
				Key:       gameRecord.Key,
				Error: &pb.Error{
					ErrorCode:    2,
					ErrorMessage: "not accessible yet",
				},
			}, nil
		}
	}

	return &pb.GameRecordValidationResult{IsSuccess: true, Key: gameRecord.Key}, nil
}

func (s *CloudsaveValidatorServer) AfterBulkReadGameRecord(ctx context.Context, gameRecords *pb.BulkGameRecord) (*pb.BulkGameRecordValidationResult, error) {
	result := []*pb.GameRecordValidationResult{}
	for _, gameRecord := range gameRecords.GetGameRecords() {
		if strings.HasSuffix(gameRecord.Key, "daily_msg") {
			var r DailyMessage
			err := json.Unmarshal(gameRecord.GetPayload(), &r)
			if err != nil {
				return nil, err
			}
			if time.Now().Before(r.AvailableOn) {
				result = append(result, &pb.GameRecordValidationResult{
					IsSuccess: false,
					Key:       gameRecord.Key,
					Error: &pb.Error{
						ErrorCode:    2,
						ErrorMessage: "not accessible yet",
					},
				})
			} else {
				result = append(result, &pb.GameRecordValidationResult{
					IsSuccess: true,
					Key:       gameRecord.Key,
				})
			}
		} else {
			result = append(result, &pb.GameRecordValidationResult{
				IsSuccess: true,
				Key:       gameRecord.Key,
			})
		}
	}

	return &pb.BulkGameRecordValidationResult{ValidationResults: result}, nil
}

func (s *CloudsaveValidatorServer) BeforeWritePlayerRecord(ctx context.Context, request *pb.PlayerRecord) (*pb.PlayerRecordValidationResult, error) {
	if strings.HasSuffix(request.GetKey(), "favourite_weapon") {
		var r CustomPlayerRecord
		err := json.Unmarshal(request.GetPayload(), &r)
		if err != nil {
			return nil, err
		}
		if err = r.Validate(); err != nil {
			errorDetail := &pb.Error{
				ErrorCode:    1,
				ErrorMessage: err.Error(),
			}

			return &pb.PlayerRecordValidationResult{
				IsSuccess: false,
				Key:       request.Key,
				UserId:    request.UserId,
				Error:     errorDetail,
			}, nil
		}
	}

	return &pb.PlayerRecordValidationResult{IsSuccess: true, Key: request.Key, UserId: request.UserId}, nil
}

func (s *CloudsaveValidatorServer) AfterReadPlayerRecord(ctx context.Context, playerRecord *pb.PlayerRecord) (*pb.PlayerRecordValidationResult, error) {
	return &pb.PlayerRecordValidationResult{IsSuccess: true, Key: playerRecord.Key, UserId: playerRecord.UserId}, nil
}

func (s *CloudsaveValidatorServer) AfterBulkReadPlayerRecord(ctx context.Context, playerRecords *pb.BulkPlayerRecord) (*pb.BulkPlayerRecordValidationResult, error) {
	result := []*pb.PlayerRecordValidationResult{}

	for _, record := range playerRecords.GetPlayerRecords() {
		result = append(result, &pb.PlayerRecordValidationResult{IsSuccess: true, Key: record.Key, UserId: record.UserId})
	}

	return &pb.BulkPlayerRecordValidationResult{ValidationResults: result}, nil
}

func (s *CloudsaveValidatorServer) BeforeWriteAdminGameRecord(ctx context.Context, request *pb.AdminGameRecord) (*pb.GameRecordValidationResult, error) {
	if strings.HasSuffix(request.GetKey(), "map") {
		var r CustomGameRecord
		err := json.Unmarshal(request.GetPayload(), &r)
		if err != nil {
			return nil, err
		}
		if err = r.Validate(); err != nil {
			errorDetail := &pb.Error{
				ErrorCode:    1,
				ErrorMessage: err.Error(),
			}

			return &pb.GameRecordValidationResult{
				IsSuccess: false,
				Key:       request.Key,
				Error:     errorDetail,
			}, nil
		}
	}

	return &pb.GameRecordValidationResult{IsSuccess: true, Key: request.Key}, nil
}

func (s *CloudsaveValidatorServer) BeforeWriteAdminPlayerRecord(ctx context.Context, request *pb.AdminPlayerRecord) (*pb.PlayerRecordValidationResult, error) {
	if strings.HasSuffix(request.GetKey(), "player_activity") {
		var r PlayerActivity
		err := json.Unmarshal(request.GetPayload(), &r)
		if err != nil {
			return nil, err
		}
		if err = r.Validate(); err != nil {
			errorDetail := &pb.Error{
				ErrorCode:    1,
				ErrorMessage: err.Error(),
			}

			return &pb.PlayerRecordValidationResult{
				IsSuccess: false,
				Key:       request.Key,
				UserId:    request.UserId,
				Error:     errorDetail,
			}, nil
		}
	}

	return &pb.PlayerRecordValidationResult{IsSuccess: true, Key: request.Key, UserId: request.UserId}, nil
}

func NewCloudsaveValidationServiceServer() *CloudsaveValidatorServer {
	return &CloudsaveValidatorServer{}
}
