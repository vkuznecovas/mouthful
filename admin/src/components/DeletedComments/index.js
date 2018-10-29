import { h, Component } from 'preact';
import Thread from '../Thread';
import Comment from '../Comment';

export default class DeletedComments extends Component {
  render() {
    // NOTE: this looks like shit. gonna rewrite it to make it more readable

    const threads = this.props.threads.sort((a, b) => {
      a = new Date(a.CreatedAt);
      b = new Date(b.CreatedAt);
      return a > b ? -1 : a < b ? 1 : 0;
    }).map(t => {
      const comments = this.props.comments.filter(c => c.ThreadId === t.Id)
        .filter(c => c.DeletedAt !== null);

      const commentsComponent = comments.map(c => <Comment id={c.Id} author={c.Author} body={c.Body} updateCommentsState={this.props.updateCommentsState} hardDelete="true" undoDelete="true" />);

      return <Thread>{commentsComponent}</Thread>;
    });

    return (
      <div>
        <h1>DELETED COMMENTS</h1>
        {threads}
      </div>
    );
  }
}