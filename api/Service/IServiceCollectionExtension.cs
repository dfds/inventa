using Microsoft.Extensions.DependencyInjection;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using Microsoft.Extensions.Configuration;

namespace Service
{
    public static class IServiceCollectionExtension
    {
        public static IServiceCollection AddServiceProxyServiceCollection(this IServiceCollection services, IConfigurationSection conf)
        {
            services.Configure<ServiceProxySettings>(
                options =>
                {
                    options.proxyUrl = conf.GetSection("Urls").GetChildren().Select(x => x.Value).ToArray();
                    options.clientId = conf.GetSection("ClientId").Value;
                    options.clientSecret = conf.GetSection("ClientSecret").Value;
                    options.clientScopes = conf.GetSection("ClientScopes").Value;
                });

            services.AddTransient<IServiceProxyService, ServiceProxyService>();
            return services;
        }
    }
}
