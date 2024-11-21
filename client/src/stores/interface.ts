export interface BaseSuccessResponse<T> {
    status: string
    message: string
    data: T
}

export interface BaseFailedResponse {
    status: string
    message: string
    error: string
}

export interface LoginResponseDto {
    token: string
}

export interface UserInfoResponseDto {
    id: number
    email: string
}