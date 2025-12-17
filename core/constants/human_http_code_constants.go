package constants

import "backend/core/types"

const (
	InternalServerError      types.HumanErrorCode = "errors.http.internal_server_error"
	BadRequestError          types.HumanErrorCode = "errors.http.bad_request_error"
	UnauthorizedError        types.HumanErrorCode = "errors.http.unauthorized"
	ForbiddenError           types.HumanErrorCode = "errors.http.forbidden"
	NotFoundError            types.HumanErrorCode = "errors.http.not_found"
	TooManyRequestsError     types.HumanErrorCode = "errors.http.too_many_requests"
	UnprocessableEntityError types.HumanErrorCode = "errors.http.unprocessable_entity"
)
