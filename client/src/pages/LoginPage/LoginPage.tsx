import { Button, Input } from "antd"
import { observer } from "mobx-react-lite"
import { ChangeEvent, useEffect } from "react"
import { authStore } from "../../stores/AuthStore"
import { InputTitle, LoginPageContent, LoginPageWrapper } from "./styles"
import { reaction } from "mobx"
import { useNavigate } from "react-router-dom"

export default observer(function () {
    const navigate = useNavigate();

	function handleSubmit() {
		console.log("loginpage submitLogin")
		authStore.submitLogin()
	}

	function handleChangeEmail(event: ChangeEvent<HTMLInputElement>) {
		authStore.email = event.currentTarget.value
	}

	function handleChangePassword(event: ChangeEvent<HTMLInputElement>) {
		authStore.password = event.currentTarget.value
	}

	useEffect(() => {
		const disposer = reaction(
			() => authStore.isAuthenticated,
			newValue => {
				if (newValue === true) {
					navigate("/success")
				}
			}
		)

		return () => disposer()
	}, [navigate])

	return (
		<LoginPageWrapper>
			<LoginPageContent>
				<InputTitle>Email</InputTitle>
				<Input onChange={handleChangeEmail} />
				<InputTitle>Пароль</InputTitle>
				<Input type="password" onChange={handleChangePassword} />
				<Button type="primary" loading={authStore.isLoading} onClick={handleSubmit}>
					Войти
				</Button>
			</LoginPageContent>
		</LoginPageWrapper>
	)
})
