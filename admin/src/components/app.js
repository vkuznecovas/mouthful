import { h, Component } from 'preact';
import Header from './header';
import Panel from '../routes/panel';

if (module.hot) {
	require('preact/debug');
}

export default class App extends Component {
	render() {
		return (
			<div id="app">
				<Header />
				<Panel />
			</div>
		);
	}
}
