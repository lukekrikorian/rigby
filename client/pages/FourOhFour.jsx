import React, { Component } from "react";

import Header from "../components/Header";

export default class FourOhFour extends Component {
	render() {
		return (
			<div>
				<Header/>
				<h2>404! That page was not found.</h2>
			</div>
		)
	}
}