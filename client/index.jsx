import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import React from "react";
import ReactDOM from "react-dom";

import Home from "./pages/Home";
import FourOhFour from "./pages/FourOhFour";
import Post from "./pages/Post";
import Signup from "./pages/Signup";

import "./css/index.css";

ReactDOM.render(
	<Router>
		<Switch>
			<Route exact path="/" component={Home}/>
			<Route path="/signup" component={Signup}/>
			<Route path="/posts/:id" component={Post}/>
			<Route component={FourOhFour}/>
		</Switch>
	</Router>,
	document.getElementById('app')
);