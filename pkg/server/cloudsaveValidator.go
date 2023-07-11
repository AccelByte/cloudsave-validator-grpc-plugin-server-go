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

func (s *CloudsaveValidatorServer) BeforeWriteGameRecord(ctx context.Context, request *pb.GameRecord) (*pb.ValidationResult, error) {
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

			return &pb.ValidationResult{
				IsSuccess: false,
				Key:       request.Key,
				Error:     errorDetail,
			}, nil
		}
	}

	return &pb.ValidationResult{IsSuccess: true, Key: request.Key}, nil
}

func (s *CloudsaveValidatorServer) AfterReadGameRecord(ctx context.Context, gameRecord *pb.GameRecord) (*pb.ValidationResult, error) {
	if strings.HasSuffix(gameRecord.Key, "daily_msg") {
		var r DailyMessage
		err := json.Unmarshal(gameRecord.GetPayload(), &r)
		if err != nil {
			return nil, err
		}
		if time.Now().Before(r.AvailableOn) {
			return &pb.ValidationResult{
				IsSuccess: false,
				Key:       gameRecord.Key,
				Error: &pb.Error{
					ErrorCode:    2,
					ErrorMessage: "not accessible yet",
				},
			}, nil
		}
	}

	return &pb.ValidationResult{IsSuccess: true, Key: gameRecord.Key}, nil
}

func (s *CloudsaveValidatorServer) AfterBulkReadGameRecord(ctx context.Context, gameRecords *pb.BulkGameRecord) (*pb.BulkValidationResult, error) {
	result := []*pb.ValidationResult{}
	for _, gameRecord := range gameRecords.GetGameRecords() {
		if strings.HasSuffix(gameRecord.Key, "daily_msg") {
			var r DailyMessage
			err := json.Unmarshal(gameRecord.GetPayload(), &r)
			if err != nil {
				return nil, err
			}
			if time.Now().Before(r.AvailableOn) {
				result = append(result, &pb.ValidationResult{
					IsSuccess: false,
					Key:       gameRecord.Key,
					Error: &pb.Error{
						ErrorCode:    2,
						ErrorMessage: "not accessible yet",
					},
				})
			} else {
				result = append(result, &pb.ValidationResult{
					IsSuccess: true,
					Key:       gameRecord.Key,
				})
			}
		}
		result = append(result, &pb.ValidationResult{
			IsSuccess: true,
			Key:       gameRecord.Key,
		})
	}

	return &pb.BulkValidationResult{ValidationResults: result}, nil
}

func (s *CloudsaveValidatorServer) BeforeWritePlayerRecord(ctx context.Context, request *pb.PlayerRecord) (*pb.ValidationResult, error) {
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

			return &pb.ValidationResult{
				IsSuccess: false,
				Key:       request.Key,
				Error:     errorDetail,
			}, nil
		}
	}

	return &pb.ValidationResult{IsSuccess: true, Key: request.Key}, nil
}

func (s *CloudsaveValidatorServer) AfterReadPlayerRecord(ctx context.Context, playerRecord *pb.PlayerRecord) (*pb.ValidationResult, error) {
	return &pb.ValidationResult{IsSuccess: true, Key: playerRecord.Key}, nil
}

func (s *CloudsaveValidatorServer) AfterBulkReadPlayerRecord(ctx context.Context, playerRecords *pb.BulkPlayerRecord) (*pb.BulkValidationResult, error) {
	result := []*pb.ValidationResult{}

	for _, record := range playerRecords.GetPlayerRecords() {
		result = append(result, &pb.ValidationResult{IsSuccess: true, Key: record.Key})
	}

	return &pb.BulkValidationResult{ValidationResults: result}, nil
}

func (s *CloudsaveValidatorServer) BeforeWriteAdminGameRecord(ctx context.Context, request *pb.AdminGameRecord) (*pb.ValidationResult, error) {
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

			return &pb.ValidationResult{
				IsSuccess: false,
				Key:       request.Key,
				Error:     errorDetail,
			}, nil
		}
	}

	return &pb.ValidationResult{IsSuccess: true, Key: request.Key}, nil
}

func (s *CloudsaveValidatorServer) AfterReadAdminGameRecord(ctx context.Context, gameRecord *pb.AdminGameRecord) (*pb.ValidationResult, error) {
	if strings.HasSuffix(gameRecord.Key, "daily_msg") {
		var r DailyMessage
		err := json.Unmarshal(gameRecord.GetPayload(), &r)
		if err != nil {
			return nil, err
		}
		if time.Now().Before(r.AvailableOn) {
			return &pb.ValidationResult{
				IsSuccess: false,
				Key:       gameRecord.Key,
				Error: &pb.Error{
					ErrorCode:    2,
					ErrorMessage: "not accessible yet",
				},
			}, nil
		}
	}

	return &pb.ValidationResult{IsSuccess: true, Key: gameRecord.Key}, nil
}

func (s *CloudsaveValidatorServer) AfterBulkReadAdminGameRecord(ctx context.Context, gameRecords *pb.BulkAdminGameRecord) (*pb.BulkValidationResult, error) {
	result := []*pb.ValidationResult{}
	for _, gameRecord := range gameRecords.GetAdminGameRecords() {
		if strings.HasSuffix(gameRecord.Key, "daily_msg") {
			var r DailyMessage
			err := json.Unmarshal(gameRecord.GetPayload(), &r)
			if err != nil {
				return nil, err
			}
			if time.Now().Before(r.AvailableOn) {
				result = append(result, &pb.ValidationResult{
					IsSuccess: false,
					Key:       gameRecord.Key,
					Error: &pb.Error{
						ErrorCode:    2,
						ErrorMessage: "not accessible yet",
					},
				})
			} else {
				result = append(result, &pb.ValidationResult{
					IsSuccess: true,
					Key:       gameRecord.Key,
				})
			}
		} else {
			result = append(result, &pb.ValidationResult{
				IsSuccess: true,
				Key:       gameRecord.Key,
			})
		}
	}

	return &pb.BulkValidationResult{ValidationResults: result}, nil
}

func (s *CloudsaveValidatorServer) BeforeWriteAdminPlayerRecord(ctx context.Context, request *pb.AdminPlayerRecord) (*pb.ValidationResult, error) {
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

			return &pb.ValidationResult{
				IsSuccess: false,
				Key:       request.Key,
				Error:     errorDetail,
			}, nil
		}
	}

	return &pb.ValidationResult{IsSuccess: true, Key: request.Key}, nil
}

func (s *CloudsaveValidatorServer) AfterReadAdminPlayerRecord(ctx context.Context, playerRecord *pb.AdminPlayerRecord) (*pb.ValidationResult, error) {
	return &pb.ValidationResult{IsSuccess: true, Key: playerRecord.Key}, nil
}

func (s *CloudsaveValidatorServer) AfterBulkReadAdminPlayerRecord(ctx context.Context, gameRecords *pb.BulkAdminPlayerRecord) (*pb.BulkValidationResult, error) {
	result := []*pb.ValidationResult{}
	for _, record := range gameRecords.GetAdminPlayerRecords() {
		result = append(result, &pb.ValidationResult{IsSuccess: true, Key: record.Key})
	}

	return &pb.BulkValidationResult{ValidationResults: result}, nil
}

func NewCloudsaveValidationServiceServer() *CloudsaveValidatorServer {
	return &CloudsaveValidatorServer{}
}
