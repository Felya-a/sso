import { styled } from "styled-components"

export const AuthWrapper = styled.div`
	display: grid;
	height: 100%;

	place-items: center;
`
export const AuthContent = styled.div`
	display: grid;
	place-items: center;
	height: min-content;

	height: 500px;
	width: 400px;

	border: 3px solid rgb(20 47 98);
	border-radius: 15px;
    background-color: rgb(15 27 50);

	grid-template-rows: 4fr 1fr;
`

// export const BackButton = styled.div`
// 	position: absolute;
// 	left: 40px;
// 	top: 40px;
// 	padding: 10px;
// 	border-radius: 10px;
// 	svg {
// 		width: 50px;
// 		height: 50px;
// 		transform: rotate(-90deg);
// 	}

// 	cursor: pointer;
// `

export const Logo = styled.div`
	font-size: 24px;
`