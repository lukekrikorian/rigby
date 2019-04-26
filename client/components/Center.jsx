import React, { Component } from "react";


export default class Center extends Component {
	constructor(props){
		super(props);
		this.getProp = this.getProp.bind(this);
	}

	getProp(prop){
		return this.props[prop] === true ? prop : "";
	}

	render(){

		const classes = `wrap ${this.getProp("vertical")} ${this.getProp("full-width")}`;

		return (
			<div className={classes}>
				{ this.props.children }
			</div>
		)
	}
}