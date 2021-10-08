using k8s.Models;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;
using Service.Classes;
using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Text.Json;
using System.Threading.Tasks;
using Microsoft.Identity.Client;

namespace Service
{
    public class ServiceProxy : IServiceProxy
    {
        private readonly HttpClient _client = new HttpClient();
        private readonly IConfidentialClientApplication _confidentialClient;
        private readonly string _scopes;
        private readonly bool _useAuth;

        public ServiceProxy(string proxyUrl, string scopes, bool useAuth, IConfidentialClientApplication confidentialClient)
        {
            _client.BaseAddress = new Uri(proxyUrl);
            _confidentialClient = confidentialClient;
            _scopes = scopes;
            _useAuth = useAuth;
        }

        public async Task<ServiceProxyResult> GetResults()
        {
            ServiceProxyResult result = new ServiceProxyResult(_client.BaseAddress.ToString());
            AuthenticationResult tokenResult = null;


            tokenResult = await _confidentialClient.AcquireTokenForClient(new List<string>(new[] { _scopes })).ExecuteAsync();


            using (_client)
            {     
                _client.BaseAddress = new Uri(_client.BaseAddress.ToString());

                _client.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", tokenResult.AccessToken);


                try
                {
                    var temp = await _client.GetAsync("/api/get-all").Result.Content.ReadAsStringAsync();
                    //TODO: Go with one deserialisation/parsing library. It'll do for now.
                    // Currently the issue is that Extensionsv1beta1Ingress & V1APIService uses IntstrIntOrString which at the moment can't be deserialised by System.Text.Json but only Newtonsoft.Json. Using CrossplaneResources as JToken doesn't work when our API controller uses System.Text.Json rather than Newtonsoft.Json. A custom converter of sorts might do the trick.
                    // See: https://docs.microsoft.com/en-us/dotnet/standard/serialization/system-text-json-how-to?pivots=dotnet-5-0#deserialization-behavior
                    var json = JObject.Parse(temp);
                    var jsonDocument = JsonDocument.Parse(temp);
                    var ingress = json["Ingress"];
                    var service = json["Service"];
                    var crossplaneResources = jsonDocument.RootElement.GetProperty("CrossplaneResources");

                    result.crossplaneResources = crossplaneResources;

                    foreach (var x in ingress)
                    {
                        var s = JsonConvert.DeserializeObject<Extensionsv1beta1Ingress>(x.ToString());
                        result.ingresses.Add(s);
                    }

                    foreach (var y in service)
                    {
                        var s = JsonConvert.DeserializeObject<V1APIService>(y.ToString());
                        result.services.Add(s);
                    }
                }
                catch (Exception e)
                {
                    Console.WriteLine(e);
                }
            }

            return result;
        }
    }
}
