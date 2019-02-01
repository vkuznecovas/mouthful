import { Component } from 'preact';
import Router, { route } from 'preact-router';
import axios from 'axios';
import Header from './Header';
import AllComments from './AllComments';
import PendingComments from './PendingComments';
import DeletedComments from './DeletedComments';
import Login from './Login';

if (module.hot) {
  require('preact/debug');
}

export default class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      config: [],
      threads: [],
      comments: [],
      authorized: false,
    };

    this.fetchConfig = this.fetchConfig.bind(this);
    this.handleLogin = this.handleLogin.bind(this);
    this.loadComments = this.loadComments.bind(this);
    this.loadThreads = this.loadThreads.bind(this);
    this.updateCommentsState = this.updateCommentsState.bind(this);
  }

  async fetchConfig() {
    try {
      // const res = await axios.get(`${root.window.location.origin}/v1/admin/config`);
      const res = await axios.get(`${window.location.origin}/v1/admin/config`);
      const config = JSON.parse(res);

      this.setState({
        configLoaded: true,
        config,
      });

    } catch (err) {
      this.setState({
        configLoaded: true,
        error: true,
      });

      console.log('Error while fetching config', err);
    }
  }

  handleLogin() {
    this.setState({ authorized: true });
  }

  updateCommentsState(comments) {
    this.setState({ comments });
  }

  async loadComments() {
    // NOTE: loading comments are implemented in "main" component, because there's only
    // one route to get them all at once. When appropriate API calls will be
    // implemented, loading certain type of comments will be moved to their components

    try {
      const comments = await axios.get(`${window.location.origin}/v1/admin/comments/all`);
      this.setState({ comments });
    } catch (err) {
      console.log('Error while loading comments', err);
    }
  }

  async loadThreads() {
    try {
      const threads = await axios.get(`${window.location.origin}/v1/admin/threads`);
      this.setState({ threads });
    } catch (err) {
      console.log('Error while loading threads', err);
    }
  }

  componentWillMount() {
    this.fetchConfig();
  }

  componentDidMount() {
    if (this.state.authorized !== true) {
      route('/login');
    }
  }

  componentDidUpdate() {
    if (this.state.authorized === true) {
      route('/');
      this.loadThreads();
      this.loadComments();
    }
  }

  render() {
    return (
      <div id="app">
        <Header />
        <div>
          <Router>
            <AllComments path="/" config={this.state.config} comments={this.state.comments} threads={this.state.threads} updateCommentsState={this.updateCommentsState} />
            <Login path="/login" config={this.state.config} handleLogin={this.handleLogin} />
            <PendingComments path="/pending" config={this.state.config} comments={this.state.comments} threads={this.state.threads} updateCommentsState={this.updateCommentsState}  />
            <DeletedComments path="/deleted" config={this.state.config} comments={this.state.comments} threads={this.state.threads} updateCommentsState={this.updateCommentsState}  />
          </Router>
        </div>
      </div>
    );
  }
}