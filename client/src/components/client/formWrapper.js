import { h, Component } from "preact";
import style from "./style";
import timeago from "./timeago"
import cookies from "./cookies"
import config from "./config"
import Form from "./form"

function getStyle(c) {
    return config.useDefaultStyle ? style[c] : c
}


export default class FormWrapper extends Component {
    constructor(props) {
        super(props);
    }

    render(props) {
        if (!this.props.comment.Confirmed) {
            return null
        }
        return <div>
            <input
                class={getStyle("mouthful_reply_button")}
                onClick={() => this.props.flipFormVisibility(this.props.comment.Id)}
                type="submit"
                value={this.props.visible ? "Close" : "Reply"}>
            </input>
            <Form id={this.props.comment.Id} visible={this.props.visible} author={this.props.author} comment={""} replyTo={this.props.comment.Id} submitForm={this.props.submitForm} />
        </div>

    }

}
