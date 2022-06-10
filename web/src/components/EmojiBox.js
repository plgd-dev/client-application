import PropTypes from "prop-types";
import React, { Component } from "react";
import { CopyToClipboard } from "react-copy-to-clipboard";
class EmojiBox extends Component {
	static propTypes = {
		title: PropTypes.string,
		Symbol: PropTypes.string,
	};

	state = {
		copied: false,
		text1: "Click to Copy",
	};

	render() {
		const HexCode = this.props.symbol.codePointAt(0).toString(16);
		const src = `//cdn.jsdelivr.net/emojione/assets/png/${HexCode}.png`;
		return (
			<CopyToClipboard
				text={this.props.symbol}
				onCopy={() =>
					this.setState({
						copied: true,
						text1: "Copied",
					})
				}
			>
				<div className="emojiBox">
					<p> {this.state.text1} </p>
					<img alt={this.props.title} src={src} />
					<span className="title"> {this.props.title} </span>
				</div>
			</CopyToClipboard>
		);
	}
}

export default EmojiBox;
