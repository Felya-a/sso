import { Button, Input } from "antd"
import { observer } from "mobx-react-lite"
import { ChangeEvent, useEffect } from "react"
import { authStore } from "../../stores/AuthStore"
import { InputTitle, RegisterPageContent, RegisterPageWrapper } from "./styles"
import { reaction } from "mobx"
import { useNavigate } from "react-router-dom"
import { Logo } from "../../layouts/AuthLayout/styles"

export default observer(function () {
    const navigate = useNavigate();

	function handleSubmit() {
		console.log("Registerpage submitRegister")
		authStore.submitRegister()
	}

	function handleChangeEmail(event: ChangeEvent<HTMLInputElement>) {
		authStore.email = event.currentTarget.value
	}

	function handleChangePassword(event: ChangeEvent<HTMLInputElement>) {
		authStore.password = event.currentTarget.value
	}

	// useEffect(() => {
	// 	if (authStore.isAuthenticated) {
	// 		navigate("/success")
	// 	}
	// }, [authStore.isAuthenticated])

	return (
		<RegisterPageWrapper>
			<Logo>
				Регистрация
			</Logo>
			<RegisterPageContent>
				<InputTitle>Имя</InputTitle>
				<Input value={authStore.nickname!} onChange={handleChangeEmail} />
				<InputTitle>Email</InputTitle>
				<Input value={authStore.email!} onChange={handleChangeEmail} />
				<InputTitle>Пароль</InputTitle>
				<Input type="password" value={authStore.password!} onChange={handleChangePassword} />
				<Button type="primary" loading={authStore.isLoading} onClick={handleSubmit}>
					Войти
				</Button>
			</RegisterPageContent>
		</RegisterPageWrapper>
	)
})
