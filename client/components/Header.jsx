import React, { Component } from "react";
import { Link } from "react-router-dom";

export default class Header extends Component {
	render () {
		return (
			<div id="headerWrap">
				<div id="headerBorder"></div>
				<Link id="header" to="/">RIGBY</Link>
			</div>
		)
	}
}