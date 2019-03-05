import React, { Component } from "react";
import { Redirect } from "react-router-dom";
import Header from "../components/Header";
import Center from "../components/Center";
import { setLoggedIn } from "../misc/LoggedIn";

class Login extends Component {
	constructor(props){
		super(props);
		this.state = {
			username: '',
			password: '',
			saveSession: false,
			error: '',
			redirect: false,
		};
		this.handleChange = this.handleChange.bind(this);
		this.submit = this.submit.bind(this);
	}

	handleChange(event){
		const { target } = event;
		if (target.type === "checkbox") { 
			this.setState({ [target.name]: target.checked });
		} else {
			this.setState({ [target.name]: target.value });
		}
	}

	submit(event){
		event.preventDefault();

		const body = {
			username: this.state.username, 
			password: this.state.password,
			saveSession: this.state.saveSession
		}

		fetch("/api/login", {
			method: "POST",
			credentials: "same-origin",
			body: JSON.stringify(body)
		}).then(response => {
			if (!response.ok) {
				response.text().then(error => this.setState({ error }));
				return
			}
			setLoggedIn(true);
			this.setState({ redirect: true });
		}).catch(console.error);
	}

	render(){
		if (this.state.redirect) {
			return <Redirect to="/"/>
		}

		return (
			<div>
				<Header/>
				<Center>
					<h1 style={{ marginBottom: 4, color: "#3c3c3c" }}>Login</h1>
					<form>
						<input type="text" onChange={this.handleChange} name="username" placeholder="Username"/>
						<br/>
						<input type="password" onChange={this.handleChange} name="password" placeholder="Password"/>
						<br/>
						<input style={{marginLeft: 0}} onChange={this.handleChange} type="checkbox" name="saveSession" id="saveSession"></input>
						<label style={{display: "inline-block", verticalAlign: "top", marginTop: 4, marginLeft: 3}} for="saveSession">Keep me logged in</label>
						<br/>
						<input type="submit" value="Submit" onClick={this.submit}/>
						<p>{this.state.error}</p>
					</form>
				</Center>
			</div>
		)
	}
}

export default Login;