import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import React from "react";
import ReactDOM from "react-dom";

import Home from "./pages/Home";
import UserHome from "./pages/UserHome";
import FourOhFour from "./pages/FourOhFour";
import Post from "./pages/Post";
import Signup from "./pages/Signup";
import Profile from "./pages/Profile";
import Logout from "./pages/Logout";
import Login from "./pages/Login";

import "./css/index.css";

window.isLoggedIn = false;
fetch("/api/isLoggedIn", { credentials: 'same-origin' })
	.then(response => {
		if (response.ok) { window.isLoggedIn = true; }
	}).catch(console.error);

ReactDOM.render(
	<Router>
		<Switch>
			<Route exact path="/" component={window.isLoggedIn ? UserHome : Home}/>
			<Route path="/logout" component={Logout}/>
			<Route path="/signup" component={Signup}/>
			<Route path="/login" component={Login}/>
			<Route path="/posts/:id" component={Post}/>
			<Route path="/~:user" component={Profile}/>
			<Route component={FourOhFour}/>
		</Switch>
	</Router>,
	document.getElementById('app')
);