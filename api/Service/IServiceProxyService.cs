using Service.Classes;
using System;
using System.Collections.Generic;
using System.Text;
using System.Threading.Tasks;

namespace Service
{
    public interface IServiceProxyService
    {
        Task<List<ServiceProxyResult>> GetResults();
    }
}
