using System;
using System.Collections.Generic;
using System.Text;

namespace Service
{
    public class ServiceProxySettings
    {
        public string[] proxyUrl { get; set; }
        public bool authEnabled { get; set; }
        public string clientId { get; set; }
        public string clientSecret { get; set; }
        public string clientScopes { get; set; }
        public string tenantId { get; set; }
    }
}
