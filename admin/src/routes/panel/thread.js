import { h, Component } from 'preact';
import style from './style';

export default class Thread extends Component {
	constructor(props) {
		super(props);
		this.state = { thread: props.thread, comments: props.comments };
	}

	render() {
		const comments = this.state.comments.map(comment => {
			return <div>{JSON.stringify(comment)}</div>;
		})
		return (
			<div class="thread">
				<div>{JSON.stringify(this.state.thread)}</div>
				<div class="comments">
					{comments}
				</div>
			</div>	
		);
	}
}
