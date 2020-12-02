exports.handler = function (event, context, callback) {
    console.log('Received event:', JSON.stringify(event, null, 2));
    var headers = event.headers;

    if (headers.Cookie) {
        callback(null, generateAllow('me', event.methodArn));
    } else {
        callback("Unauthorized");
    }
}

// Help function to generate an IAM policy
var generatePolicy = function (principalId, effect, resource) {
    var authResponse = {};
    authResponse.principalId = principalId;
    if (effect && resource) {
        var policyDocument = {};
        policyDocument.Version = '2012-10-17';
        policyDocument.Statement = [];
        var statementOne = {};
        statementOne.Action = 'execute-api:Invoke';
        statementOne.Effect = effect;
        statementOne.Resource = resource;
        policyDocument.Statement[0] = statementOne;
        authResponse.policyDocument = policyDocument;
    }
    return authResponse;
}

var generateAllow = function (principalId, resource) {
    return generatePolicy(principalId, 'Allow', resource);
}
