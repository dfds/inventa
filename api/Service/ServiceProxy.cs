using k8s;
using k8s.Models;
using Microsoft.Extensions.Options;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;
using Service.Classes;
using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Net.Http.Headers;
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
                    var json = JObject.Parse(temp);
                    var ingress = json["Ingress"];
                    var service = json["Service"];

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
