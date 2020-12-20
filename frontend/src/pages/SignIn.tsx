import React from 'react';
import {
    GoogleLogin,
    GoogleLoginResponse,
    GoogleLoginResponseOffline, 
} from 'react-google-login';

const responseGoogleSuccess = (response: GoogleLoginResponse | GoogleLoginResponseOffline) => {


    if ((response as GoogleLoginResponse).profileObj) {
        const user = response as GoogleLoginResponse;
        const gToken = user.getAuthResponse().id_token;

        fetch("https://api2.zcclock.com/sign-in", {
            method: "POST",
            body: JSON.stringify({googleToken: gToken})
        })
        .then(response => response.json())
        .then(response => {
            console.log(response);
            /*
            this.setState({
                users: response.Users,
                newUser: "",
                newUserFieldEnabled: true,
            });
            */
        })
        .catch(err => {
            // Need to indicate error...
            console.log(err);
        });

    } else {
        // consider failure?
        //this.failureResponse(response)
    }
}

const responseGoogleFailure = (response: any) => {
    console.log(response);
}

export const SignIn = () => {
    return (<div>
        <h2>Sign In</h2>
        <br />
        <GoogleLogin
            clientId={process.env.REACT_APP_GOOGLE_SIGNIN_CLIENT_ID || ""}
            buttonText="Login"
            onSuccess={responseGoogleSuccess}
            onFailure={responseGoogleFailure}
            cookiePolicy={'single_host_origin'}
        />
        <div className="g-signin2" data-onsuccess="onSignIn"></div>
    </div>);
}