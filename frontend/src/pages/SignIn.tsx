import React from "react"

import {
  GoogleOAuthProvider,
  GoogleLogin,
  CredentialResponse,
} from "@react-oauth/google"
import {
  Navigate,
  //useLocation,
} from "react-router-dom"
import { StandardFetch } from "./utils/FetchHelper"

type Props = {}

type State = {
  message: string
  messageClass: string
  processing: boolean
  successfullyLoggedIn: boolean
}

const ERROR_NOT_AUTHORIZED = "not authorized"
const ERROR_NOT_AUTHENTICATED = "not authenticated"
const ERROR_UNKNOWN = "unknown"

export class SignIn extends React.Component<Props, State> {
  state: Readonly<State> = {
    message: "",
    messageClass: "",
    processing: false,
    successfullyLoggedIn: false,
  }

  componentDidMount() {
    const state = new URLSearchParams(window.location.search).get("state")
    if (state === "401") {
      this.setAlert("", ERROR_NOT_AUTHENTICATED)
    } else if (state === "403") {
      this.setAlert("", ERROR_NOT_AUTHORIZED)
    }
  }

  responseGoogleSuccess(response: CredentialResponse) {
    this.setAlert("", "")

    if (response) {
      const user = response as CredentialResponse
      const gToken = user.credential

      this.setState({ processing: true })
      StandardFetch("sign-in", {
        method: "POST",
        body: JSON.stringify({ googleToken: gToken }),
      })
        .then((response) => {
          switch (response.status) {
            case 403:
              this.setAlert("", ERROR_NOT_AUTHORIZED)
              return
            case 200:
              this.setState({ successfullyLoggedIn: true })
              return
            default:
              console.log("unhandled response code: " + response.status)
              this.setAlert("", ERROR_UNKNOWN)
              return
          }
        })
        .catch((err) => {
          console.log(err)
          this.setAlert("", ERROR_UNKNOWN)
        })
        .finally(() => {
          this.setState({ processing: false })
        })
    } else {
      this.setAlert("", ERROR_UNKNOWN)
    }
  }

  responseGoogleFailure(response: any) {
    this.setAlert("", ERROR_UNKNOWN)
    console.log(response)
  }

  signOut() {
    this.setAlert("", "")
    this.setState({ processing: true })
    StandardFetch("sign-in", { method: "DELETE" })
      .then((response) => {
        this.setAlert("Signed out successfully", "")
      })
      .catch((err) => {
        console.log(err)
        this.setAlert("", ERROR_UNKNOWN)
      })
      .finally(() => {
        this.setState({ processing: false })
      })
  }

  setAlert(message: string, errorCode: string) {
    if (message === "" && errorCode === "") {
      this.setState({
        message: "",
        messageClass: "",
      })
      return
    }

    let messageClass = "alert-primary"
    if (message === "") {
      messageClass = "alert-danger"
      switch (errorCode) {
        case "":
          // Do nothing.
          break
        case ERROR_NOT_AUTHORIZED:
          message = "Not Authorized - request access from an administrator"
          break
        case ERROR_NOT_AUTHENTICATED:
          message = "Not Authenticated - log in to authenticate"
          break
        default: // Includes `ERROR_UNKNOWN`.
          message = "An error has occurred"
      }
    }

    this.setState({
      message: message,
      messageClass: messageClass,
    })
  }

  render() {
    if (this.state.successfullyLoggedIn) {
      const next = new URLSearchParams(window.location.search).get("next")
      if (next != null && next !== "") {
        return <Navigate to={next} replace />
      }
      return <Navigate to="/devices/" replace />
    }

    let innerContent = this.renderNonProcessing()
    if (this.state.processing) {
      innerContent = this.renderProcessing()
    }

    return (
      <div>
        <h2>Sign In</h2>
        <br />
        {innerContent}
      </div>
    )
  }

  renderProcessing() {
    return <div>Processing...</div>
  }

  renderNonProcessing() {
    return (
      <GoogleOAuthProvider
        clientId={process.env.REACT_APP_GOOGLE_SIGNIN_CLIENT_ID || ""}
      >
        <div>
          {this.state.message && (
            <div className={"alert " + this.state.messageClass} role="alert">
              {this.state.message}
            </div>
          )}
          <br />
          <GoogleLogin
            onSuccess={this.responseGoogleSuccess.bind(this)}
            onError={() => {
              this.responseGoogleFailure.bind(this)
            }}
            auto_select
            useOneTap
          />
        </div>
      </GoogleOAuthProvider>
    )
  }
}
