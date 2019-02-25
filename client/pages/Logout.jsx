import React, { Component } from "react";
import { Redirect } from "react-router-dom";
import { setLoggedIn } from "../misc/LoggedIn";

class Logout extends Component {
	constructor(props){
		super(props);
		this.state = { redirect: false };
	}
	
	componentDidMount(){
		fetch("/api/logout", { method: "POST", credentials: "same-origin" })
			.then(response => {
				if (response.ok) {
					setLoggedIn(false);
					this.setState({ redirect: true });
				}
			}).catch(console.error);
	}

	render(){
		if (this.state.redirect) {
			return <Redirect to="/"/>
		}
	}
}

export default Logout;