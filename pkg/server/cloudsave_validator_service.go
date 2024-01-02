// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	pb "cloudsave-validator-grpc-plugin-server-go/pkg/pb"
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

func (s *CloudsaveValidatorServer) BeforeWriteGameBinaryRecord(ctx context.Context, request *pb.GameBinaryRecord) (*pb.GameRecordValidationResult, error) {
	if strings.HasSuffix(request.GetKey(), "event_banner") && request.GetBinaryInfo() != nil {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, request.GetBinaryInfo().GetUrl(), nil)
		if err != nil {
			return nil, err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		fileSize, err := strconv.Atoi(resp.Header.Get("Content-Length"))
		if err != nil {
			return nil, err
		}

		if fileSize/1000 > MaxSizeEventBannerInKB {
			errorDetail := &pb.Error{
				ErrorCode:    1,
				ErrorMessage: fmt.Sprintf("maximum size for event banner is %d kB", MaxSizeEventBannerInKB),
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

func (s *CloudsaveValidatorServer) AfterReadGameBinaryRecord(ctx context.Context, request *pb.GameBinaryRecord) (*pb.GameRecordValidationResult, error) {
	if strings.HasSuffix(request.GetKey(), "daily_event_stage") && request.GetBinaryInfo() != nil {
		if !isSameDate(time.Now().UTC(), request.GetBinaryInfo().GetUpdatedAt().AsTime().UTC()) {
			errorDetail := &pb.Error{
				ErrorCode:    1,
				ErrorMessage: fmt.Sprintf("today's %s is not ready yet", request.Key),
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

func (s *CloudsaveValidatorServer) AfterBulkReadGameBinaryRecord(ctx context.Context, request *pb.BulkGameBinaryRecord) (*pb.BulkGameRecordValidationResult, error) {
	result := []*pb.GameRecordValidationResult{}

	for _, record := range request.GetGameBinaryRecords() {
		if strings.HasSuffix(record.GetKey(), "daily_event_stage") && record.GetBinaryInfo() != nil {
			if !isSameDate(time.Now().UTC(), record.GetBinaryInfo().GetUpdatedAt().AsTime().UTC()) {
				errorDetail := &pb.Error{
					ErrorCode:    1,
					ErrorMessage: fmt.Sprintf("today's %s is not ready yet", record.Key),
				}

				result = append(result, &pb.GameRecordValidationResult{
					IsSuccess: false,
					Key:       record.Key,
					Error:     errorDetail,
				})
			} else {
				result = append(result, &pb.GameRecordValidationResult{IsSuccess: true, Key: record.Key})
			}
		}
	}

	return &pb.BulkGameRecordValidationResult{ValidationResults: result}, nil
}

func (s *CloudsaveValidatorServer) BeforeWritePlayerBinaryRecord(ctx context.Context, request *pb.PlayerBinaryRecord) (*pb.PlayerRecordValidationResult, error) {
	if strings.HasSuffix(request.GetKey(), "id_card") && request.GetBinaryInfo() != nil {
		if request.GetBinaryInfo().GetVersion() > 1 {
			errorDetail := &pb.Error{
				ErrorCode:    1,
				ErrorMessage: "id card can only be created once",
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

func (s *CloudsaveValidatorServer) AfterReadPlayerBinaryRecord(ctx context.Context, request *pb.PlayerBinaryRecord) (*pb.PlayerRecordValidationResult, error) {
	return &pb.PlayerRecordValidationResult{IsSuccess: true, Key: request.Key, UserId: request.UserId}, nil
}

func (s *CloudsaveValidatorServer) AfterBulkReadPlayerBinaryRecord(ctx context.Context, request *pb.BulkPlayerBinaryRecord) (*pb.BulkPlayerRecordValidationResult, error) {
	result := []*pb.PlayerRecordValidationResult{}

	for _, record := range request.GetPlayerBinaryRecords() {
		result = append(result, &pb.PlayerRecordValidationResult{IsSuccess: true, Key: record.Key, UserId: record.UserId})
	}

	return &pb.BulkPlayerRecordValidationResult{ValidationResults: result}, nil
}

func NewCloudsaveValidationServiceServer() *CloudsaveValidatorServer {
	return &CloudsaveValidatorServer{}
}
