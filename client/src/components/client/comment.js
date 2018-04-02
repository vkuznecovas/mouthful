import { h, Component } from "preact";
import style from "./style";
import timeago from "./timeago"
import cookies from "./cookies"
import Form from "./form"


function formatDate(d) {
    var dd = new Date(d)
    return timeago(dd)
}


export default class Comment extends Component {
    constructor(props) {
        super(props);
        this.getStyle = this.getStyle.bind(this);
    }
    getStyle(c) {
        return this.props.config.useDefaultStyle ? style[c] : c
    }
    render(props) {
        return <div>
        <div class={this.getStyle("mouthful_author")}>{this.props.comment.Author}
        <span class={this.getStyle("mouthful_date")}>{formatDate(this.props.comment.CreatedAt)}</span>
        {(!this.props.comment.Confirmed && this.props.config.moderation) ? <span class={this.getStyle("mouthful_moderation")}>In queue for moderation</span> : null}
        </div>
        <div class={this.getStyle("mouthful_comment_body")} dangerouslySetInnerHTML={{ __html: this.props.comment.Body }} />
        </div>
    }

}
