import { h, Component } from 'preact';
import axios from 'axios';
import style from './style';

export default class Login extends Component {
  constructor(props) {
    super(props);
    this.state = {
      value: '',
    };

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleOauthClick = this.handleOauthClick.bind(this);
  }

  handleOauthClick(provider) {
    window.location.replace(`${this.props.url}v1/oauth/auth/${provider}`);
  }

  handleChange(event) {
    this.setState({ value: event.target.value });
  }

  async handleSubmit(event) {
    event.preventDefault();

    const url = `${window.location.origin}/v1/admin/login`;
    const config = {
      headers: {
        'Content-Type': 'application/json',
      },
    };

    try {
      const res = await axios.post(url, JSON.stringify({ password: this.state.value }), config);

      if (res.status === 204) {
        this.props.handleLogin();
      }
    } catch (err) {
      console.log('Something went wrong', err);
    }

  }

  componentDidMount() {
    if (this.props.config.disablePasswordLogin) {
      this.props.handleLogin();
    }
  }

  render() {
    const providersListItems = this.props.config.oauthProviders && this.props.config.oauthProviders.length > 0
      ? this.props.config.oauthProviders.map( x => <li><a onClick={() => this.handleOauthClick(x)}>Log in with {x}</a></li>)
      : null;

    const oauthProviders = <div><ul>{providersListItems}</ul></div>;

    return (
      <div class={style.mouthful_login}>
        <form onSubmit={this.handleSubmit}>
          <input type="password" value={this.state.value} onChange={this.handleChange} />
          <input class={style.mouthful_submit} type="submit" value="Submit" />
        </form>
        {oauthProviders}
      </div>
    );
  }
}