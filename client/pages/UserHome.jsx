import React, { Component } from "react";
import Header from "../components/Header";
import Center from "../components/Center";

class UserHome extends Component {
	render(){
		return (
			<div>
				<Header/>
				<Center>
					<h1>Welcome back to rigby B)</h1>
				</Center>
			</div>
		);
	}
}

export default UserHome;