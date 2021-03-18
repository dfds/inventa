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
        public static IServiceCollection AddServiceProxyServiceCollection(this IServiceCollection services, IConfiguration conf)
        {

            services.Configure<ServiceProxySettings>(
                options =>
                {
                    options.authEnabled = bool.Parse(conf.GetSection("INVENTA_API_AUTH_ENABLE").Value);
                    options.proxyUrl = conf.GetSection("INVENTA_API_OPERATOR_URLS").Value.Split(',');
                    options.clientId = conf.GetSection("INVENTA_API_AUTH_CLIENT_ID").Value;
                    options.clientSecret = conf.GetSection("INVENTA_API_AUTH_CLIENT_SECRET").Value;
                    options.clientScopes = conf.GetSection("INVENTA_API_AUTH_CLIENT_SCOPES").Value;
                    options.tenantId = conf.GetSection("INVENTA_API_AUTH_TENANT_ID").Value;
                });

            services.AddTransient<IServiceProxyService, ServiceProxyService>();
            return services;
        }
    }
}