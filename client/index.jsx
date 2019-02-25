import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import React from "react";
import ReactDOM from "react-dom";

import Home from "./pages/Home";
import FourOhFour from "./pages/FourOhFour";
import Post from "./pages/Post";
import Signup from "./pages/Signup";
import Profile from "./pages/Profile";
import Logout from "./pages/Logout";
import Login from "./pages/Login";
import Browse from "./pages/Browse";
import CreatePost from "./pages/CreatePost";
import Conversation from "./pages/Conversation";

import { setLoggedIn } from "./misc/LoggedIn";
import "./css/index.css";

fetch("/api/isLoggedIn", { credentials: 'same-origin' })
	.then(response => {
		setLoggedIn(response.ok);
	}).catch(console.error);

ReactDOM.render(
	<Router>
		<Switch>
			<Route exact path="/" component={Home}/>
			<Route path="/logout" component={Logout}/>
			<Route path="/signup" component={Signup}/>
			<Route path="/login" component={Login}/>
			<Route path="/posts/:id" component={Post}/>
			<Route path="/post" component={CreatePost}/>
			<Route path="/~:user" component={Profile}/>
			<Route path="/browse/:page(popular|recent)" component={Browse}/>
			<Route path="/conversation" component={Conversation}/>
			<Route component={FourOhFour}/>
		</Switch>
	</Router>,
	document.getElementById('app')
);