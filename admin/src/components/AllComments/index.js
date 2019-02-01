import { Component } from 'preact';
import Thread from '../Thread';
import Comment from '../Comment';

export default class AllComments extends Component {
  render() {
    // TODO: rewrite it to make it more readible

    const threads = this.props.threads.sort((a, b) => {
      a = new Date(a.CreatedAt);
      b = new Date(b.CreatedAt);

      if (a > b) {
        return -1;
      } else if (a < b) {
        return 1;
      } else {
        return 0;
      }
    }).map(t => {
      const comments = this.props.comments.filter(c => c.ThreadId === t.Id)
        .filter(c => c.DeletedAt === null);

      const commentsComponent = comments.map(c => <Comment id={c.Id} author={c.Author} body={c.Body} updateCommentsState={this.props.updateCommentsState} />);

      return <Thread>{commentsComponent}</Thread>;
    });

    return (
      <div>
        <h1>ALL COMMENTS</h1>
        {threads}
      </div>
    );
  }
}