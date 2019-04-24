import React, { Component } from "react";

class Highlight extends Component {
	render(){
		return (
			<div class="highlight">
				<h1>{ this.props.children }</h1>
			</div>
		)
	}
}

export default Highlight;