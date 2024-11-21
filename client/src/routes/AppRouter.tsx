import { Route, Routes } from "react-router-dom"
import AuthLayout from "../layouts/AuthLayout"
import LoginPage from "../pages/LoginPage/LoginPage"
import Page404 from "../pages/Page404/Page404"

const AppRouter = () => {
	return (
		<Routes>
			<Route
				path="/login"
				element={
					<AuthLayout>
						<LoginPage />
					</AuthLayout>
				}
			/>

			<Route
				path="/"
				element={
					// <>
					// 	<div>Userid: {authStore.userid}</div>
					// 	<div>Email: {authStore.email}</div>
					// 	<div>Token: {authStore.token}</div>
					// 	<div>Path=/</div>
					// 	<button onClick={handleButton}>button</button>
					// </>
					<AuthLayout>
						<LoginPage />
					</AuthLayout>
				}
			/>

			<Route path="*" element={<Page404 />} />
		</Routes>
	)
}

export default AppRouter
