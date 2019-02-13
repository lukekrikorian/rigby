import React, { Component } from "react";
import { Link } from "react-router-dom";

import "../css/Header.css";

export default class Header extends Component {
	render () {
		return (
			<div className="headerWrap">
				<div id="headerBorder"></div>
				<Link id="header" to="/">🐾 rigby react</Link>
			</div>
		)
	}
}