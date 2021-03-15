using k8s.Models;
using System;
using System.Collections.Generic;
using System.Text;

namespace Service.Classes
{
    public class ServiceProxyResult
    {
        public string proxyName { get; set; }
        public List<Extensionsv1beta1Ingress> ingresses { get; set; }
        public List<V1APIService> services { get; set; }

        public ServiceProxyResult(string name)
        {
            proxyName = name;
            ingresses = new List<Extensionsv1beta1Ingress>();
            services = new List<V1APIService>();
        }
    }
}
