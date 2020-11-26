import React from 'react';
import {
    GoogleLogin,
    GoogleLoginResponse,
    GoogleLoginResponseOffline, 
} from 'react-google-login';

const responseGoogleSuccess = (response: GoogleLoginResponse | GoogleLoginResponseOffline) => {
    console.log(response);
}

const responseGoogleFailure = (response: any) => {
    console.log(response);
}

export const SignIn = () => {
    return (<div>
        <h2>Sign In</h2>
        <br />
        <GoogleLogin
            clientId="278545785364-vpm7qvrccmultq5rml71auhq5qa7co97.apps.googleusercontent.com"
            buttonText="Login"
            onSuccess={responseGoogleSuccess}
            onFailure={responseGoogleFailure}
            cookiePolicy={'single_host_origin'}
        />
        <div className="g-signin2" data-onsuccess="onSignIn"></div>
    </div>);
}