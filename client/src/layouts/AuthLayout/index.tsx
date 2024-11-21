import { FC } from "react"
import { useNavigate } from "react-router-dom"
import { AuthContent, AuthWrapper, Logo } from "./styles"

interface AuthLayoutProps {
	children: React.ReactNode
}

const AuthLayout: FC<AuthLayoutProps> = ({ children }) => {
	const navigate = useNavigate()

	return (
		<AuthWrapper>
			<AuthContent>
				<Logo>
					Авторизация
				</Logo>
				{children}
			</AuthContent>
		</AuthWrapper>
	)
}

export default AuthLayout
