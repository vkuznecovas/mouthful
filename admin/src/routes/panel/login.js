import { h, Component } from 'preact';
import style from './style';


export default class Login extends Component {
	constructor(props) {
        super(props);
		this.state = { value: '' };
        this.onLogin = props.onLogin;
        
		this.handleChange = this.handleChange.bind(this);
		this.handleSubmit = this.handleSubmit.bind(this);
		this.handleOauthClick = this.handleOauthClick.bind(this);
	}
	handleOauthClick(provider) {
		window.location.replace(this.props.url + "v1/oauth/auth/" + provider);
	}
	handleChange(event) {
		this.setState({ value: event.target.value });
	}

	handleSubmit(event) {
		var context = this;
        
		event.preventDefault();		
		var http = new XMLHttpRequest();
		var url = this.props.url + "v1/admin/login";
		http.open("POST", url, true);
		
		//Send the proper header information along with the request
		http.setRequestHeader("Content-type", "application/json");
		http.onreadystatechange = function() {//Call a function when the state changes.
			if(http.readyState == 4 && http.status == 204) {
                context.onLogin();
			} 
		}
		http.send(JSON.stringify({password: context.state.value}));
	}
	render() {
		var login = <form onSubmit={this.handleSubmit}>
		<label class={style.passwordTitle}>Password:</label>
		<input type="password" value={this.state.value} onChange={this.handleChange} />
		<input class={style.mouthful_submit}type="submit" value="Submit" />
		</form>
		var loginDiv = this.props.config.disablePasswordLogin ? null : login;
		var providersListItems = this.props.config.oauthProviders && this.props.config.oauthProviders.length > 0 
		? this.props.config.oauthProviders.map(x => <li class={style.mouthful_admin_li}><a class={style.mouthful_admin_oauth_a} onClick={() => this.handleOauthClick(x)}>Log in with {x}</a></li>)
		: null
		var oauthProviders = <div>
			<ul>
				{providersListItems}
			</ul>
		</div>
		return (
			<div class={style.mouthful_login}>
				{loginDiv}
				{oauthProviders}
			</div>
		);
	}
}
