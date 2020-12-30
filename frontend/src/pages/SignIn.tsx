import React from 'react';
import {
    GoogleLogin,
    GoogleLoginResponse,
    GoogleLoginResponseOffline, 
} from 'react-google-login';

type Props = {};

type State = {
    error: string;
};

const ERROR_NOT_AUTHORIZED = 'not authorized';
const ERROR_UNKNOWN = 'unknown';

export class SignIn extends React.Component<Props, State> {
    state: Readonly<State> = {
        error: "",
    }

    responseGoogleSuccess(response: GoogleLoginResponse | GoogleLoginResponseOffline) {
        if ((response as GoogleLoginResponse).profileObj) {
            const user = response as GoogleLoginResponse;
            const gToken = user.getAuthResponse().id_token;

            fetch((process.env.REACT_APP_BACKEND_DOMAIN || "") + "/sign-in", {
                method: "POST",
                credentials: "include",
                body: JSON.stringify({googleToken: gToken})
            })
            .then(response => {
                switch(response.status) {
                    case 403:
                        this.setState({ error: ERROR_NOT_AUTHORIZED });
                        return
                }

                // this is a promise: console.log(response.json());
                // grab token from response and redirect to somewhere useful?
            }).catch(err => {
                console.log(err);
                this.setState({ error: ERROR_UNKNOWN });
            });
        } else {
            // consider failure?
            //this.failureResponse(response)
        }
    }

    responseGoogleFailure(response: any) {
        console.log(response);
    }

    render() {
        let errorMessage = ""
        switch(this.state.error) {
            case "":
                // Do nothing.
                break;
            case ERROR_NOT_AUTHORIZED:
                errorMessage = "Not Authorized - request access from an administrator"
                break;
            default: // Includes `ERROR_UNKNOWN`.
                errorMessage = "An error has occurred"
        }
        return (<div>
            <h2>Sign In</h2>
            <br />
            {
                errorMessage &&
                <div className="alert alert-danger" role="alert">{errorMessage}</div>
            }
            <br />
            <GoogleLogin
                clientId={process.env.REACT_APP_GOOGLE_SIGNIN_CLIENT_ID || ""}
                buttonText="Login"
                onSuccess={this.responseGoogleSuccess.bind(this)}
                onFailure={this.responseGoogleFailure.bind(this)}
                cookiePolicy={'single_host_origin'}
            />
            <div className="g-signin2" data-onsuccess="onSignIn"></div>
        </div>);
    }
}
