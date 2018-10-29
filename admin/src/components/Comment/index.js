import { h, Component } from 'preact';
import axios from 'axios';
import Button from '../Button';

export default class Comment extends Component {
  constructor(props) {
    super(props);
    this.state = {
      author: this.props.author,
      body: this.props.body,
    };

    this.updateComment = this.updateComment.bind(this);
    this.deleteComment = this.deleteComment.bind(this);
    this.undoDeleteComment = this.undoDeleteComment.bind(this);
    this.confirmComment = this.confirmComment.bind(this);
    this.handleAuthorChange = this.handleAuthorChange.bind(this);
    this.handleBodyChange = this.handleBodyChange.bind(this);
  }

  handleAuthorChange(event) {
    this.setState({ author: event.target.value });
  }

  handleBodyChange(event) {
    this.setState({ body: event.target.value });
  }

  async updateComment() {
    try {
      const data = {
        CommentId: this.props.id,
        Body: this.props.body,
        Author: this.props.author
      };
      const res = await axios.patch(`${window.location.origin}/v1/admin/comments`, data);

      if (res.status === 204) {
        // TODO: reload comments
        console.log('Hooray! Comment updating successful');
      }
    } catch (err) {
      console.log('Error while updating comment', err);
    }
  }

  async confirmComment() {
    // NOTE: dunno what should I think about this
    try {
      const data = {
        // CommentId: this.props.id,
        // Body: this.props.body,
        // Author: this.props.author,
        Confirmed: true,
      };
      const res = await axios.patch(`${window.location.origin}/v1/admin/comments`, data);

      if (res.status === 204) {
        // TODO: reload comments
        console.log('Hooray! Comment updating successful');
      }
    } catch (err) {
      console.log("Error while confirming comment", err);
    }
  }

  // TODO: confirm (only on PendingComments)
  async deleteComment(hardDelete) {
    try {
      const config = hardDelete
        ? {
          data: {
            CommentId: this.props.id,
            Hard: true,
          }
        }
        : {
          data: {
            CommentId: this.props.id,
          }
        };

      const res = await axios.delete(`${window.location.origin}/v1/admin/comments`, config);

      if (res.status === 204) {
        const comments = this.props.comments.filter(c => c.Id !== this.props.id);
        this.props.updateCommentsState(comments);
      }
    } catch (err) {
      console.log('Error while deleting comment', err);
    }
  }

  async undoDeleteComment() {
    try {
      const data = { CommentId: this.props.id };
      const res = await axios.post(`${window.location.origin}/v1/admin/comments/restore`, data);

      if (res.status === 204) {
        const comments = this.props.comments.filter(c => c.Id !== this.props.id);
        this.props.updateCommentsState(comments);
        // update
      }
    } catch (err) {
      console.log('Error while restoring comment', err);
    }
  }

  render() {
    const deleteButton = this.props.hardDelete
      ? <Button text="Hard delete" onClick={this.deleteComment(true)} />
      : <Button text="Delete" onClick={this.deleteComment} />;

    const undoDeleteButton = this.props.undoDelete
      ? <Button text="Undo delete" onClick={this.undoDeleteComment} />
      : null;

    const confirmCommentButton = this.props.confirmComment
      ? <Button text="Confirm" onClick={this.confirmComment} />
      : null;

    return (
      <div>
        <h1>{this.props.author}</h1>
        <input value={this.state.author} onChange={this.handleAuthorChange} />
        <p>{this.props.body}</p>
        <textarea value={this.state.body} onChange={this.handleBodyChange}></textarea>
        <div>
          <Button text="Update" onClick={this.updateComment} />
          {confirmCommentButton}
          {undoDeleteButton}
          {deleteButton}
        </div>
      </div>
    );
  }
}