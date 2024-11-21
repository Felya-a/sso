import { notification } from "antd"
import axios from "axios"
import { makeAutoObservable } from "mobx"
import { BaseFailedResponse, BaseSuccessResponse, LoginResponseDto, UserInfoResponseDto } from "./interface"

class AuthStore {
	token: string | null = null
	userid: number | null = null
	email: string | null = null
	password: string | null = null
	errors: string[] = []
	isAuthenticated = false
	isLoading = false

	constructor() {
		makeAutoObservable(this)
	}

	async login(email: string, password: string) {
		try {
			const loginResponse = await axios.post<BaseSuccessResponse<LoginResponseDto>>(process.env.REACT_APP_API_URL + "/login", {
				email,
				password
			})

			this.token = loginResponse.data.data.token

			const userInfo = await this.userInfo(this.token)

			if (!userInfo) {
				throw new Error("not found user info")
			}

			this.userid = userInfo.data.id
			this.email = userInfo.data.email

			this.isAuthenticated = true
		} catch (error) {
			console.error("Login failed:", error)
			this.isAuthenticated = false
			if (axios.isAxiosError<BaseFailedResponse>(error)) {
				if (error.response?.data.error) {
					this.setError(error.response?.data.error)
					return
				}
				this.setError("Неверные логин или пароль")
				return
			}
			this.setError("Ошибка")
			return
		}
	}

	async userInfo(token: string): Promise<BaseSuccessResponse<UserInfoResponseDto> | undefined> {
		try {
			const userInfoReponse = await axios.get<BaseSuccessResponse<UserInfoResponseDto>>(process.env.REACT_APP_API_URL + "/userinfo", {
				headers: {
					Authorization: `Bearer ${token}`
				}
			})

			return userInfoReponse.data
		} catch (error) {
			console.error("Get user info failed:", error)
		}
	}

	logout() {
		this.token = null
		this.isAuthenticated = false
		console.log("Logged out")
	}

	setError(message: string) {
		notification.error({ message, duration: 3 })
	}

	async submitLogin() {
		console.log("authstore submitLogin")
		this.isLoading = true

		if (!this.email) {
			this.setError("Введите email")
			this.isLoading = false
			return
		}

		if (!this.password) {
			this.setError("Введите пароль")
			this.isLoading = false
			return
		}

		await this.login(this.email, this.password)
		this.isLoading = false
	}
}

export const authStore = new AuthStore()
