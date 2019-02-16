import React, { Component } from "react";
import "../css/Heart.css";

class Heart extends Component {
	constructor(props){
		super(props);
		this.state = { error: false };
		this.sendLike = this.sendLike.bind(this);
	}

	sendLike(e){
		e.preventDefault();
		fetch(`/api/vote/${this.props.post}`, { 
			method: "POST",
			credentials: "same-origin"
		}).then(response => {
				if (!response.ok) {
					this.setState({ error: true });
				}
		}).catch(console.error)
	}

	render(){
		return (
			<a href="#" className={`${this.state.error ? 'heartError' : ''} heart`} onclick={this.sendLike}>
				<img src="/static/images/heart.svg" alt="Heart button"/>
			</a>
		)
	}
}

export default Heart;
