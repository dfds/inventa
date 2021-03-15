using Microsoft.AspNetCore.Components.WebAssembly.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Text;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Components.WebAssembly.Authentication;

namespace DFDSServiceUi
{
    public class Program
    {
        public static async Task Main(string[] args)
        {
            var builder = WebAssemblyHostBuilder.CreateDefault(args);

            builder.Services.AddMsalAuthentication(options =>
            {
                builder.Configuration.Bind("AzureAd", options.ProviderOptions.Authentication);
                options.ProviderOptions.Cache.CacheLocation = "localStorage";
                options.ProviderOptions.DefaultAccessTokenScopes.Add(
                    builder.Configuration["AzureAd:ClientScopes"]);
                options.UserOptions.RoleClaim = "roles";
                options.ProviderOptions.LoginMode = "redirect";

            });

            builder.Services.AddScoped<CustomAuthorizationMessageHandler>();

            builder.Services.AddHttpClient("DFDSServiceApi", client => client.BaseAddress = new Uri(builder.Configuration["ApiUrl"]))
                .AddHttpMessageHandler<CustomAuthorizationMessageHandler>();
            //builder.Services.AddScoped(sp => sp.GetRequiredService<IHttpClientFactory>().CreateClient("DFDSServiceApi"));


            builder.RootComponents.Add<App>("#app");

            //builder.Services.AddScoped(sp => new HttpClient { BaseAddress = new Uri(builder.HostEnvironment.BaseAddress) });

            await builder.Build().RunAsync();
        }
    }
}
