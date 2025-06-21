import { notification } from "antd"
import axios from "axios"
import { makeAutoObservable } from "mobx"
import { BaseFailedResponse, BaseSuccessResponse, LoginResponseDto, UserInfoResponseDto } from "./interface"

class AuthStore {
	authorizationCode: string | null = null
	token: string | null = null
	userid: number | null = null
	email: string | null = null
	password: string | null = null
	nickname: string | null = null
	// email: string | null = "11@gmail.com" // TODO: DEBUG ONLY
	// password: string | null = "12345678" // TODO: DEBUG ONLY
	errors: string[] = []
	isAuthenticated = false
	isLoading = false

	constructor() {
		makeAutoObservable(this)
	}

	async login(email: string, password: string) {
		try {
			// const url = new URL(window.location.href);
			// const params = new URLSearchParams(new URL(window.location.href).search);
			// url.search        - '?client_id=talk-client&redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Fcallback'
			// params.toString() - 'client_id=talk-client&redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Fcallback'

			const params = new URLSearchParams(new URL(window.location.href).search);
			const redirectUrl = params.get("redirect_url")

			if (!redirectUrl) {
				this.setError("Отсутсвует redirect_url в query параметрах")
				return
			}
			
			// Формируем URL с параметрами
			const baseUrl = process.env.REACT_APP_SSO_SERVER_URL + "/login";
			const url = new URL(baseUrl);
			params.forEach((value, key) => {
				url.searchParams.append(key, value);
			});

			console.log(url)
			console.log(url.toString())
			const loginResponse = await axios.post<BaseSuccessResponse<LoginResponseDto>>(url.toString(), {
				email,
				password
			})

			console.log(loginResponse)
			this.authorizationCode = loginResponse?.data?.data?.authorization_code
			if (this.authorizationCode) {
				const callbackUrl = new URL(redirectUrl)
				callbackUrl.searchParams.set("authorization_code", this.authorizationCode)
				console.log(callbackUrl.toString())

				window.location.href = callbackUrl.toString()
			} else {
				this.setError("Нет данных")
				return
			}

			// const userInfo = await this.userInfo(this.token)

			// if (!userInfo) {
			// 	throw new Error("not found user info")
			// }

			// this.userid = userInfo.data.id
			// this.email = userInfo.data.email

			// this.isAuthenticated = true
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

	// async userInfo(token: string): Promise<BaseSuccessResponse<UserInfoResponseDto> | undefined> {
	// 	try {
	// 		const userInfoReponse = await axios.get<BaseSuccessResponse<UserInfoResponseDto>>(process.env.REACT_APP_SSO_SERVER_URL + "/userinfo", {
	// 			headers: {
	// 				Authorization: `Bearer ${token}`
	// 			}
	// 		})

	// 		return userInfoReponse.data
	// 	} catch (error) {
	// 		console.error("Get user info failed:", error)
	// 	}
	// }

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

	async submitRegister() {
		console.log("authstore submitRegister")
		this.isLoading = true

		if (!this.nickname) {
			this.setError("Введите имя")
			this.isLoading = false
			return
		}

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

		// TODO: implemented
		// await this.register(this.email, this.password)
		this.isLoading = false
	}
}

export const authStore = new AuthStore()
