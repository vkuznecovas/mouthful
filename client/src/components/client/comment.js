import { h, Component } from "preact";
import style from "./style";
import timeago from "./timeago"
import cookies from "./cookies"
import config from "./config"
import Form from "./form"

function getStyle(c) {
    return config.useDefaultStyle ? style[c] : c
}

function formatDate(d) {
    var dd = new Date(d)
    return timeago(dd)
}


export default class Comment extends Component {
    constructor(props) {
        super(props);
    }
   
    render(props) {
        return <div>
        <div class={getStyle("mouthful_author")}>{this.props.comment.Author}
        <span class={getStyle("mouthful_date")}>{formatDate(this.props.comment.CreatedAt)}</span>
        {(!this.props.comment.Confirmed && config.moderation) ? <span class={getStyle("mouthful_moderation")}>In queue for moderation</span> : null}
        </div>
        <div class={getStyle("mouthful_comment_body")} dangerouslySetInnerHTML={{ __html: this.props.comment.Body }} />
        </div>
    }

}
