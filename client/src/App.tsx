import "./App.css"
import { BrowserRouter as Router } from "react-router-dom"
import AppRouter from "./routes/AppRouter"
import { ConfigProvider } from "antd"

function App() {
	return (
		<ConfigProvider
			theme={{
				token: {
					fontFamily: `'Lora', sans-serif` // Задайте ваш шрифт
				}
			}}
		>
			<Router>
				<AppRouter />
			</Router>
		</ConfigProvider>
	)
}

export default App
