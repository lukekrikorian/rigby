import React, { Component } from "react";
import "../css/Center.css";

export default class Center extends Component {
	constructor(props){
		super(props);
	}
	render(){
		return (
			<div className="centerWrap">
				<div className="center">
					{ this.props.children }
				</div>
			</div>
		)
	}
}