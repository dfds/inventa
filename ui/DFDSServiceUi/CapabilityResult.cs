using k8s.Models;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace DFDSServiceUi
{
    public class CapabilityResult
    {
        public string capabilityName { get; set; }
        public List<Extensionsv1beta1Ingress> ingresses { get; set; }

        public List<V1APIService> services { get; set; }

        public CapabilityResult()
        {
            ingresses = new List<Extensionsv1beta1Ingress>();
            services = new List<V1APIService>();
        }
    }
}
