import React, { Component } from "react";
import { Link } from "react-router-dom";
import Header from "../components/Header";
import Center from "../components/Center";

class Home extends Component {
	render () {
		return (
			<div>
				<Header/>
				<Center>
					<h2 style={{ marginBottom: 0, color: "#332f2f" }}>Rigby</h2>
					<p style={{ marginTop: 5, marginBottom: 0 }}>a site where cool kids can discuss</p>
					<p style={{ marginTop: 4 }}><Link to="/signup">Sign up</Link> or <Link to="/login">Log in</Link> or <Link to="/posts/cf3ae5d0-6c27-4633-8948-c31c56f6fccb">Browse recent posts</Link></p>
				</Center>
			</div>
		)
	}
}

export default Home;