import { Navigate, Route, Routes } from "react-router-dom"
import AuthLayout from "../layouts/AuthLayout"
import LoginPage from "../pages/LoginPage/LoginPage"
import Page404 from "../pages/Page404/Page404"
import RegisterPage from "../pages/RegisterPage"

const AppRouter = () => {
	return (
		<Routes>
			<Route
				path="/"
				element={
					<AuthLayout>
						<LoginPage />
					</AuthLayout>
				}
			/>

			<Route
				path="/register"
				element={
					<AuthLayout>
						<RegisterPage />
					</AuthLayout>
				}
			/>

			{/* <Route
				path="/"
				element={
					<AuthLayout>
						<LoginPage />
					</AuthLayout>
				}
			/> */}

			<Route path="*" element={<Page404 />} />
		</Routes>
	)
}

export default AppRouter
